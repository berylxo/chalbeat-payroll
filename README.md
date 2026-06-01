# Chalbeat Payroll

Deployment notes

- Frontend (Vite/React): set `VITE_API_BASE` to your API base URL when building or running in non-local environments. If left empty, requests will use relative paths (recommended when the frontend and backend are served from the same origin).

  Example (.env.production):
  VITE_API_BASE=https://api.example.com

  To build the frontend for production (from `/frontend`):
  ```bash
  # install deps
  npm install
  # build with Vite using .env.production
  npm run build
  ```

- Backend (Go): configure allowed CORS origin via `CORS_ALLOWED_ORIGIN`. Default is `*` (development). For production, set the origin or use your frontend host.

  Example:
  ```bash
  export CORS_ALLOWED_ORIGIN=https://app.example.com
  go run ./backend
  ```

- PDF generation: ensure `wkhtmltopdf` is installed and available in PATH on the server if payslip PDF generation is required.

- Quick dev run (local):
  - Start the backend (from repository root):
    ```bash
    cd backend
    go run main.go
    ```
  - Start the frontend dev server (from repository root):
    ```bash
    cd frontend
    npm install
    npm run dev
    ```

If you want, I can also add a small shell script or Dockerfiles for containerized deployment. Let me know which target environment you plan to deploy to (Docker, Heroku, Vercel, Netlify, or a VPS) and I will scaffold deployment configs.