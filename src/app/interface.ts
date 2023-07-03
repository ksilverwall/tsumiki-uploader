export interface BackendInterface {
  upload(zipData: ArrayBuffer): Promise<string>;
}
