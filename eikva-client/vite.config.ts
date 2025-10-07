import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
    build: {
        emptyOutDir: true,
        outDir: "../static"
    },
    plugins: [react()],
    server: {
        open: true,
        proxy: {
            "/auth": "http://localhost:3000",
            "/groups": "http://localhost:3000",
            "/test-cases": "http://localhost:3000",
            "/steps": "http://localhost:3000",
            "/uploads": "http://localhost:3000",
            "/ws": {
                target: "ws://localhost:3000",
                ws: true,
                changeOrigin: true,
                secure: false
            }
        },
    },
});
