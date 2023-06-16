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

type UploadAction = {
  type: "UPLOAD";
  promises: { [key: GroupId]: Promise };
};

type Action = LoadAction | CreateGroupAction | UploadAction;

type UploadGroupAction = {
  type: "UPLOAD";
  promise: Promise<void>;
};

type GroupAction = UploadGroupAction;
