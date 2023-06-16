type LoadAction = {
  type: "LOAD";
  groupId: GroupId;
  items: { id: ItemId; file: File }[];
};

type CreateGroupAction = {
  type: "CREATE_GROUP";
  newGroupId: GroupId;
  source?: {
    groupId: GroupId;
    selected: ItemId[];
  };
};

type Action = LoadAction | CreateGroupAction;
