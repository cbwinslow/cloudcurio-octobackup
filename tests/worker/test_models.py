import pytest
from pydantic_models import User, Script, ReviewJob
from datetime import datetime

def test_user_model():
    user = User(id='1', name='cbwinslow', email='blaine.winslow@gmail.com', role='admin')
    assert user.id == '1'
    assert user.role == 'admin'
    assert user.model_dump() == {
        'id': '1',
        'name': 'cbwinslow',
        'email': 'blaine.winslow@gmail.com',
        'role': 'admin'
    }

    # Validation error
    with pytest.raises(ValueError):
        User(id='', name='test', email='invalid')

def test_script_model():
    script = Script(slug='hello.sh', content='echo hello', channel='stable')
    assert script.slug == 'hello.sh'
    assert script.downloads == 0

def test_review_job_model():
    job = ReviewJob(id='job1', repo_url='https://github.com/test/repo', status='queued')
    assert job.status == 'queued'
    assert 'created_at' in job.model_dump()
