import { BackendApi } from "./BackendApi";
import { S3SignedAccessor } from "./S3SignedAccessor";

type Accessor = {
  backendApi: BackendApi;
  s3Accessor: S3SignedAccessor;
}

export default class Backend {
  constructor(private readonly accessor: Accessor) { }
  async upload(zipData: ArrayBuffer): Promise<string> {
    const t = await this.accessor.backendApi.createTransaction();
    await this.accessor.s3Accessor.put(t.url, zipData);

    return t.id;
  }
}
