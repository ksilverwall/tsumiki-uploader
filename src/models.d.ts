type GroupId = string;
type ItemId = string;

type GroupState = "EDITING" | "ARCHIVING" | "COMPLETE";

type Item = {
  file: File;
};

type Group = {
  label: string;
  items: { [key: ItemId]: Item };
} & ({
  state: "EDITING" | "ARCHIVING";
} | {
  state: "COMPLETE";
  key: string;
});

type Status = {
  groups: { [key: GroupId]: Group };
};
