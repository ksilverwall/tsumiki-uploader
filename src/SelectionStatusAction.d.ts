type LoadAction = {
  type: "LOAD";
  groupId: GroupId;
  items: { id: ItemId; file: File }[];
};

type ArchiveAction = {
  type: "ARCHIVE";
  newGroupId: GroupId;
  groupId: GroupId;
  source: {
    selected: ItemId[];
  };
};

type CreateGroupAction = {
  type: "CREATE_GROUP";
  newGroupId: GroupId;
};

type Action = LoadAction | ArchiveAction | CreateGroupAction;
