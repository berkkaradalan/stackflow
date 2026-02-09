from dotenv import load_dotenv
import os

load_dotenv()

class Settings:
    def __init__(self):
        self.BACKEND_HOSTNAME = self._get_env("BACKEND_HOSTNAME")
        self.BACKEND_PORT = self._get_env("BACKEND_PORT")
        self.BACKEND_URL = "http://" + self.BACKEND_HOSTNAME + ":" + self.BACKEND_PORT
        

    @staticmethod
    def _get_env(var_name: str) -> str:
        value = os.getenv(var_name)
        if value is None:
            raise EnvironmentError(f"Required environment variable '{var_name}' is not set.")
        return value

settings = Settings()