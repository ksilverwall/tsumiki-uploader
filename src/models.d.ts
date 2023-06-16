type GroupId = string;
type ItemId = string;

type GroupState = "EDITING" | "ARCHIVING" | "COMPLETE";

type Item = {
  file: File;
};

type Group = {
  label: string;
  state: GroupState;
  items: { [key: ItemId]: Item };
};

type Status = {
  groups: { [key: GroupId]: Group };
};
