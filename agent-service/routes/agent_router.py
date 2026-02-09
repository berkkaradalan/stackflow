from fastapi import APIRouter, Header
import requests
from core.config import settings
import os
from app.agents.pm_agent import init_pm_agent


agent_router = APIRouter()

@agent_router.post("/agents/pm/breakdown-task")
async def breakdown_task(
    agent_id: int,
    project_id: int,
    task_title: str,
    task_description: str,
    authorization: str = Header(None)
):
    """
    Break down a task into micro-tasks using the PM agent and write them to the Go backend.
    """
    agent_data = requests.get(
        f"{settings.BACKEND_URL}/api/agents/{agent_id}",
        headers={"Authorization": authorization}
    ).json()
    
    agent = init_pm_agent(
        agent_name=agent_data["name"],
        agent_description=agent_data["description"],
        provider_name=agent_data["provider"],
        model_id=agent_data["model"],
        api_key=agent_data["api_key"]
    )
    
    os.environ["AUTH_TOKEN"] = authorization
    os.environ["PROJECT_ID"] = str(project_id)
    
    prompt = f"""
    Break down this task into micro-tasks:
    
    Title: {task_title}
    Description: {task_description}
    
    Create specific, actionable tasks using create_task_tool. 
    Each task should be small enough to be completed independently.
    Use appropriate priorities and tags.
    """
    
    # result = agent.run(prompt)

    result = await agent.ainvoke({
        "messages": [{"role": "user", "content": prompt}]
    })
    
    
    os.environ.pop("AUTH_TOKEN", None)
    os.environ.pop("PROJECT_ID", None)
    
    return {
        "status": "completed",
        "original_task": {
            "title": task_title,
            "description": task_description
        },
        "result": result
    }