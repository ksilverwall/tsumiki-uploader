import { DownloadInfo, Transaction } from "../../gen";
import { ApplicationError } from "./ApplicationError";
import axios from "axios";

export class BackendApi {
  constructor(private readonly endpointUrl: string) { }
  async createTransaction(): Promise<Transaction> {
    try {
      const response = await axios.post(`${this.endpointUrl}/storage/transactions`, {}, {
        headers: {
          'Content-Type': 'application/json'
        }
      });

      return response.data as Transaction;
    } catch (err) {
      throw new ApplicationError(`failed to createTransaction`, err);
    }
  }
  async getDownloadUrl(key: string): Promise<DownloadInfo> {
    try {
      const response = await axios.get(`${this.endpointUrl}/storage/files/${key}`, {
        headers: {
          'Content-Type': 'application/json'
        }
      });

      return response.data as DownloadInfo;
    } catch (err) {
      throw new ApplicationError(`failed to getDownloadUrl`, err);
    }
  }
}
