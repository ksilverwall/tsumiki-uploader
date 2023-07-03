import { BackendApi } from "./repositories/BackendApi";
import Backend from "./repositories/BackendSignedBufferUploader";
import { S3SignedAccessor } from "./repositories/S3SignedAccessor";

export function CreateConnector(): Backend {
  const endpoint = import.meta.env.VITE_BACKEND_API_ENDPOINT;
  if (!endpoint) {
    throw new Error("environment variable VITE_BACKEND_API_ENDPOINT is not set");
  }

  return new Backend({
    backendApi: new BackendApi(import.meta.env.VITE_BACKEND_API_ENDPOINT),
    s3Accessor: new S3SignedAccessor(),
  })
}
