import os
from deepagents import create_deep_agent
import requests
from typing import Literal
from core.env import settings
from app.providers import create_model_client

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
        agent_name: str, 
        agent_description: str,
        provider_name: str,
        model_id: str,
        api_key: str,
        temperature: float = 0.7,
    ):
    """
    Initialize a PM agent with the specified provider and model.
    
    Args:
        agent_name: Name of the agent
        agent_description: Description of the agent's role
        provider_name: Provider name (zai, anthropic, gemini, kimi, openrouter)
        model_id: Model ID from the provider registry
        api_key: API key for the provider
        temperature: Sampling temperature (default: 0.7)
    
    Returns:
        A configured Deep Agent instance
    
    Example:
        >>> agent = init_pm_agent(
        ...     agent_name="Project Manager",
        ...     agent_description="Manages software projects",
        ...     provider_name="anthropic",
        ...     model_id="claude-opus-4-5-20251101",
        ...     api_key=os.getenv("ANTHROPIC_API_KEY")
        ... )
    """
    # Create the model client using the provider registry
    model = create_model_client(
        provider_name=provider_name,
        model_id=model_id,
        api_key=api_key,
        temperature=temperature
    )
    
    pm_agent = create_deep_agent(
        agent_name=agent_name,
        agent_description=agent_description,
        model=model,
        tools=[
            list_tasks_tool,
            create_task_tool,
            analyze_project_tool
        ]
    )
    return pm_agent


def init_pm_agent_legacy(
        agent_name: str, 
        agent_description: str,
        agent_model: str,
    ):
    """
    Legacy initializer for backward compatibility.
    Uses model string directly without the provider registry.
    """
    pm_agent = create_deep_agent(
        agent_name=agent_name,
        agent_description=agent_description,
        model=agent_model,
        tools=[
            list_tasks_tool,
            create_task_tool,
            analyze_project_tool
        ]
    )
    return pm_agent
