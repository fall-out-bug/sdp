"""Tests for sdp.beads.base abstract client."""

import pytest

from sdp.beads.base import BeadsClient
from sdp.beads.models import BeadsStatus, BeadsTask, BeadsTaskCreate


def test_beads_client_is_abstract() -> None:
    """Test that BeadsClient cannot be instantiated directly."""
    with pytest.raises(TypeError, match="Can't instantiate abstract class"):
        BeadsClient()  # type: ignore


def test_beads_client_requires_create_task() -> None:
    """Test that subclass must implement create_task."""
    class IncompleteClient(BeadsClient):
        def get_task(self, task_id: str):
            pass
        def update_task_status(self, task_id: str, status: BeadsStatus):
            pass
        def get_ready_tasks(self):
            pass
        def add_dependency(self, from_id: str, to_id: str, dep_type: str = "blocks"):
            pass
        def list_tasks(self, status=None, parent_id=None):
            pass
        def update_metadata(self, task_id: str, metadata: dict):
            pass
    
    with pytest.raises(TypeError):
        IncompleteClient()  # type: ignore


def test_beads_client_requires_all_methods() -> None:
    """Test that a valid subclass implements all abstract methods."""
    class ValidClient(BeadsClient):
        def create_task(self, params: BeadsTaskCreate) -> BeadsTask:
            return BeadsTask(
                id="bd-test",
                title=params.title,
                description=params.description or "",
                status=BeadsStatus.OPEN,
            )
        
        def get_task(self, task_id: str):
            return None
        
        def update_task_status(self, task_id: str, status: BeadsStatus):
            pass
        
        def get_ready_tasks(self):
            return []
        
        def add_dependency(self, from_id: str, to_id: str, dep_type: str = "blocks"):
            pass
        
        def list_tasks(self, status=None, parent_id=None):
            return []
        
        def update_metadata(self, task_id: str, metadata: dict):
            pass
    
    # Should instantiate successfully
    client = ValidClient()
    assert isinstance(client, BeadsClient)
    
    # Test that methods can be called
    task = client.create_task(BeadsTaskCreate(title="Test", description="Test"))
    assert task.id == "bd-test"
    assert task.title == "Test"
