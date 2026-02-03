from fastapi import APIRouter
from .agent_router import agent_router
from .memoy_router import memory_router
from .llm_router import llm_router
from .code_router import code_router

main_router = APIRouter()

@main_router.get("/health")
def get_health_check():
    return {"status":"healthy", "message":"agent-service is running."}

main_router.include_router(agent_router, prefix="/agent", tags=["Agent"], responses={404: {"description": "Not found"}})
main_router.include_router(memory_router, prefix="/memory", tags=["Memory"], responses={404: {"description": "Not found"}})
main_router.include_router(llm_router, prefix="/llm", tags=["LLM"], responses={404: {"description": "Not found"}})
main_router.include_router(code_router, prefix="/code", tags=["Code"], responses={404: {"description": "Not found"}})