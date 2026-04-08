import { createRouter, createRootRoute } from "@tanstack/react-router";
import App from './App';

const rootRoute = createRootRoute({
    component: App,
});

export const router = createRouter({
    routeTree: rootRoute,
});