"""Sample Python file with PRD annotations for testing."""

from sdp.prd import prd_flow, prd_step


@prd_flow("submission-processing")
@prd_step(1, "Receive submission from queue")
async def process_submission(self, job):
    """Process a single submission."""
    pass


@prd_flow("submission-processing")
@prd_step(2, "Clone git repository")
async def clone_repository(self, url):
    """Clone the repository."""
    pass


@prd_flow("submission-processing")
@prd_step(3, "Run tests in Docker")
async def run_in_sandbox(self, path):
    """Run in isolated environment."""
    pass


@prd_flow("notification-flow")
@prd_step(1, "Send email notification")
def send_notification(self, email, result):
    """Send notification."""
    pass
