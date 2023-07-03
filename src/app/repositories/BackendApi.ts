import { Transaction } from "../../gen";
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
}
