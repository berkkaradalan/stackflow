from fastapi import FastAPI

app = FastAPI()

@app.get("/health")
def get_health_check():
    return {"status":"healthy", "message":"agent-service is running."}