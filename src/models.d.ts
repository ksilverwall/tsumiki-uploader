type GroupId = string;
type ItemId = string;

type Item = {
  file: File;
};

type Group = {
  label: string;
  state: 'EDITING' | 'ARCHIVING' | 'COMPLETE';
  items: { [key: ItemId]: Item };
};

type Status = {
  groups: { [key: GroupId]: Group };
};
