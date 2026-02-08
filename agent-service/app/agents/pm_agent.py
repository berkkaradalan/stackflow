import os
from deepagents import create_deep_agent
import requests
from typing import Literal
from core.env import settings

def analyze_project_tool(project_id: str):
    response = requests.get(f"{settings.BACKEND_URL}/api/projects/{project_id}")
    return response.json()

def create_task_tool(
    title: str,
    description: str = "",
    priority: Literal["low", "medium", "high", "critical"] = "medium",
    assigned_agent_id: int = None,
    reviewer_id: int = None,
    tags: list[str] = None
):
    payload = {
        "title": title,
        "description": description,
        "priority": priority
    }

    if assigned_agent_id:
        payload["assigned_agent_id"] = assigned_agent_id
    if reviewer_id:
        payload["reviewer_id"] = reviewer_id
    if tags:
        payload["tags"] = tags
    
    response = requests.post(f"{settings.BACKEND_URL}/api/tasks", json=payload)
    return response.json()

def list_tasks_tool(project_id: str, status: str = None):
    params = {"project_id": project_id}

    if status:
        params["status"] = status

    response = requests.get(f"{settings.BACKEND_URL}/api/tasks", params=params)

    return response.json()

def init_pm_agent(
        agent_name:str, 
        agent_description: str,
        agent_model: str,
    ):
    pm_agent = create_deep_agent(
        agent_name = agent_name,
        agent_description = agent_description,
        model=agent_model,
        tools=[
            list_tasks_tool,
            create_task_tool,
            analyze_project_tool
        ]
    )
    return pm_agent