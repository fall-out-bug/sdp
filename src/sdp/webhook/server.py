"""GitHub webhook HTTP server."""

from __future__ import annotations

import json
import logging
import os
from typing import TYPE_CHECKING, Any

if TYPE_CHECKING:
    from .handler import EventHandler

from .handler import EventHandler
from .signature import SignatureError, SignatureValidator


logger = logging.getLogger(__name__)


# Global reference to the running server instance
# Used by WebhookHandler to access the server's event handler
_server_instance: WebhookServer | None = None


def _get_server() -> WebhookServer:
    """Get the current server instance.

    Raises:
        RuntimeError: If no server is running
    """
    if _server_instance is None:
        raise RuntimeError("WebhookServer not initialized")
    return _server_instance


class WebhookHandler:  # type: ignore[misc]
    """HTTP request handler for webhooks.

    This class is instantiated by http.server.HTTPServer.
    It accesses the running WebhookServer via _get_server().
    """

    def __init__(self) -> None:
        """Initialize handler (HTTPServer calls this without args)."""
        # HTTPServer creates handler instances without arguments
        # We access the server via _get_server()
        self._server = _get_server()

    def do_POST(self) -> None:
        """Handle POST request (webhook delivery)."""
        if self.path != "/webhook":
            self.send_error(404, "Not Found")
            return

        # Get signature header
        signature = self.headers.get("X-Hub-Signature-256")

        # Read payload
        content_length = int(self.headers.get("Content-Length", 0))
        payload = self.rfile.read(content_length)

        # Handle event
        try:
            self._server.handler.handle(payload, signature)
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b'{"status": "ok"}')
        except SignatureError as e:
            logger.warning("Signature validation failed: %s", e)
            self.send_error(401, "Invalid signature")
        except Exception as e:
            logger.exception("Error handling webhook")
            self.send_error(500, str(e))

    def do_GET(self) -> None:
        """Handle GET request (health check and events)."""
        if self.path == "/health":
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b'{"status": "healthy"}')
        elif self.path == "/events":
            events = self._server.handler.get_event_log()
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
            self.wfile.write(json.dumps(events).encode())
        else:
            self.send_error(404, "Not Found")

    def log_message(self, format: str, *args: Any) -> None:  # noqa: A002
        """Suppress default HTTP logging."""
        pass


class WebhookServer:
    """HTTP server for receiving GitHub webhooks.

    Uses a simple HTTP server implementation for minimal dependencies.
    """

    def __init__(
        self,
        host: str = "0.0.0.0",
        port: int = 8080,
        webhook_secret: str | None = None,
        log_file: str = ".sdp/webhook.log",
    ) -> None:
        """Initialize webhook server.

        Args:
            host: Host to bind to
            port: Port to listen on
            webhook_secret: GitHub webhook secret for validation
            log_file: Path to event log file
        """
        self._host = host
        self._port = port
        self._running = False

        # Create signature validator and event handler
        self._validator = SignatureValidator(webhook_secret)
        self._handler = EventHandler(self._validator, log_file)

    @property
    def handler(self) -> EventHandler:
        """Get the event handler instance."""
        return self._handler

    def start(self) -> None:
        """Start the webhook server (blocking)."""
        try:
            from http.server import HTTPServer
        except ImportError:
            logger.error("http.server not available")
            return

        # Set global server reference for WebhookHandler
        global _server_instance
        _server_instance = self
        self._running = True

        # Create and bind server
        server = HTTPServer((self._host, self._port), WebhookHandler)

        logger.info("Webhook server listening on %s:%s", self._host, self._port)

        try:
            server.serve_forever()
        except KeyboardInterrupt:
            logger.info("Webhook server stopped")
        finally:
            self._running = False
            _server_instance = None
            server.server_close()

    def stop(self) -> None:
        """Stop the webhook server."""
        self._running = False


def start_server(
    host: str = "0.0.0.0",
    port: int = 8080,
    webhook_secret: str | None = None,
    smee_url: str | None = None,
) -> WebhookServer:
    """Start webhook server.

    Args:
        host: Host to bind to
        port: Port to listen on
        webhook_secret: GitHub webhook secret
        smee_url: SMEE.io URL for webhook tunneling (optional)

    Returns:
        WebhookServer instance
    """
    if smee_url:
        logger.info("SMEE tunneling URL: %s", smee_url)
        logger.info("Configure GitHub webhook to forward to this URL")

    # Get secret from env if not provided
    if not webhook_secret:
        webhook_secret = os.getenv("GITHUB_WEBHOOK_SECRET")

    server = WebhookServer(
        host=host,
        port=port,
        webhook_secret=webhook_secret,
    )

    return server


def start_server_background(
    host: str = "0.0.0.0",
    port: int = 8080,
    webhook_secret: str | None = None,
) -> WebhookServer:
    """Start webhook server in background thread.

    Args:
        host: Host to bind to
        port: Port to listen on
        webhook_secret: GitHub webhook secret

    Returns:
        WebhookServer instance
    """
    import threading

    server = WebhookServer(
        host=host,
        port=port,
        webhook_secret=webhook_secret,
    )

    thread = threading.Thread(target=server.start, daemon=True)
    thread.start()

    logger.info("Webhook server started in background on %s:%s", host, port)

    return server
