import { ApplicationError } from "./ApplicationError";
import axios from "axios";

export class S3SignedAccessor {
  async put(url: string, buffer: ArrayBuffer): Promise<void> {
    try {
      await axios.put(url, buffer, {
        headers: {
          'Content-Type': 'application/zip'
        }
      });
    } catch (err) {
      throw new ApplicationError(`failed to put data to ${url}`, err);
    }
  }
}
