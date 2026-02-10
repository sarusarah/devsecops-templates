// Frontend Application
// This is a simple example for monorepo CI/CD demonstration

class FrontendApp {
  constructor() {
    this.name = 'Frontend Application';
    this.version = '1.0.0';
  }

  greet() {
    return `Hello from ${this.name} v${this.version}`;
  }

  apiEndpoint() {
    return process.env.API_URL || 'http://localhost:8000/api';
  }
}

const app = new FrontendApp();
console.log(app.greet());
console.log(`API Endpoint: ${app.apiEndpoint()}`);

module.exports = FrontendApp;
