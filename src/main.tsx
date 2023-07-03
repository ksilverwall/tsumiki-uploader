import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import { Router as RemixRouter } from "@remix-run/router";
import UploadPage from "./components/UploadPage.tsx";
import "./index.css";
import { CreateConnector } from "./app/builder.ts";
import DownloadPage from "./components/DownloadPage.tsx";
import { Context } from "./app/context.ts";

type ApplicationConstants = {
  router: RemixRouter;
  error?: undefined;
} | {
  error: Error;
}

function init(): ApplicationConstants {
  try {
    const c: Context = {
      backend: CreateConnector()
    }
    const router = createBrowserRouter([
      {
        path: "/",
        element: <UploadPage context={c} />,
      },
      {
        path: "/download",
        element: <DownloadPage context={c} />,
      },
    ]);

    return { router: router };
  } catch (err) {
    if (err instanceof Error) {
      return { error: err };
    }
    throw err;
  }
}

const constants = init();

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    {constants.error ? <p>{constants.error.message}</p> : <RouterProvider router={constants.router} />}
  </React.StrictMode>
);
