import { Configuration, DefaultApi } from "../gen";
import Backend from "./repositories/Backend";
import { S3SignedAccessor } from "./repositories/S3SignedAccessor";

export function CreateConnector(): Backend {
  const endpoint = import.meta.env.VITE_BACKEND_API_ENDPOINT;
  if (!endpoint) {
    throw new Error("environment variable VITE_BACKEND_API_ENDPOINT is not set");
  }

  return new Backend({
    backendApi: new DefaultApi(new Configuration({ basePath: import.meta.env.VITE_BACKEND_API_ENDPOINT })),
    s3Accessor: new S3SignedAccessor(),
  })
}
