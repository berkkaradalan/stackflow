from fastapi import APIRouter

memory_router = APIRouter()

@memory_router.post("/store")
async def store_memory():
    return

@memory_router.get("/search")
async def search_memory():
    return

@memory_router.delete("/{memory_id}")
async def delete_memory(memory_id: str):
    return

@memory_router.post("/generate-embedding")
async def generate_embedding():
    return