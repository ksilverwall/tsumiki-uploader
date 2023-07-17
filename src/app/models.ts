import { ApplicationError } from "./repositories/ApplicationError";

export type GroupId = string;
export type ItemId = string;

export type GroupState = "EDITING" | "ARCHIVING" | "COMPLETE" | "FAILED";

export type Item = {
  file: File;
};

export type Group = {
  label: string;
  items: { [key: ItemId]: Item };
} & ({
  state: "EDITING" | "ARCHIVING";
} | {
  state: "COMPLETE";
  key: string;
} | {
  state: "FAILED";
  error: ApplicationError;
});

export type Status = {
  groups: { [key: GroupId]: Group };
};
