#!/bin/bash
# Sample bash script with PRD annotations for testing

# @prd: flow=deploy-flow, step=1, desc=Checkout code
git checkout main

# @prd: flow=deploy-flow, step=2, desc=Install dependencies
pip install -r requirements.txt

# @prd: flow=deploy-flow, step=3, desc=Run migrations
python manage.py migrate

# @prd: flow=deploy-flow, step=4, desc=Restart services
docker-compose restart

echo "Deploy complete"
