type MarkItemAction = {
  type: "MARK_ITEM";
  groupId: GroupId;
  index: number;
  value: boolean;
};

type LoadAction = {
  type: "LOAD";
  groupId: GroupId;
  files: File[];
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
