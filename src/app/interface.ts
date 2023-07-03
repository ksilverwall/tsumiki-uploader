export interface BackendInterface {
  upload(zipData: ArrayBuffer): Promise<string>;
  download(key: string): Promise<void>;
}
