type MarkItemAction = {
  type: "MARK_ITEM";
  groupId: GroupId;
  itemId: ItemId;
  value: boolean;
};

type LoadAction = {
  type: "LOAD";
  groupId: GroupId;
  items: { id: ItemId; file: File}[];
};

type ArchiveAction = {
  type: "ARCHIVE";
  groupId: GroupId;
};

type CreateGroupAction = {
  type: "CREATE_GROUP";
  newGroupId: GroupId;
};

type Action = MarkItemAction | LoadAction | ArchiveAction | CreateGroupAction;
