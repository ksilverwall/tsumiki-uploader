export interface BackendInterface {
  upload(zipData: ArrayBuffer): Promise<string>;
  download(key: string): Promise<void>;
  getThumbnailUrl(key: string): Promise<string[]>;
}
