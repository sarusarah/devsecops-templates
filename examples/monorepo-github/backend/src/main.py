"""
Backend Application
This is a simple example for monorepo CI/CD demonstration
"""

from typing import Dict


class BackendApp:
    """Simple backend application class"""

    def __init__(self):
        self.name = "Backend Application"
        self.version = "1.0.0"

    def greet(self) -> str:
        """Return greeting message"""
        return f"Hello from {self.name} v{self.version}"

    def health_check(self) -> Dict[str, str]:
        """Return health check status"""
        return {
            "status": "healthy",
            "service": self.name,
            "version": self.version
        }


if __name__ == "__main__":
    app = BackendApp()
    print(app.greet())
    print(f"Health: {app.health_check()}")
