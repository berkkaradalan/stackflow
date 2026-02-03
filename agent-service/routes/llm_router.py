from fastapi import APIRouter

llm_router = APIRouter()

@llm_router.post("/completion")
async def llm_completion():
    return

@llm_router.post("/chat")
async def llm_chat():
    return

@llm_router.get("/models")
async def get_models():
    return

@llm_router.get("/costs")
async def get_costs():
    return
