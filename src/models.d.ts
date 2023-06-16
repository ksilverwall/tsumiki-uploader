type GroupId = string;
type ItemId = string;

type Item = {
  file: File;
};

type Group = {
  items: { [key: ItemId]: Item };
};

type Status = {
  groups: { [key: GroupId]: Group };
};
