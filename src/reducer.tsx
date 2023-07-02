type Reducer<S, A> = (group: S, action: A) => S;

function mapDict<S>(
  groups: { [key: string]: S },
  func: (key: string, value: S) => S
): { [key: string]: S } {
  return Object.fromEntries(
    Object.keys(groups).map((groupId) => [
      groupId,
      func(groupId, groups[groupId]),
    ])
  );
}

const GroupReducer: Reducer<Group, Action> = (state, action) => {
  switch (action.type) {
    case "SET_GROUP_ITEMS":
      return {
        ...state,
        items: action.items,
      };
    case "UPLOAD":
      return {
        ...state,
        state: "ARCHIVING",
      };
    case "UPLOAD_COMPLETE":
      return {
        ...state,
        state: "COMPLETE",
        key: action.key,
      };
    default:
      return state;
  }
};

const StatusReducer: Reducer<Status, Action> = (state, action) => {
  switch (action.type) {
    case "CREATE_GROUP": {
      let groups: { [key: GroupId]: Group };
      const source = action.source;
      if (source) {
        const marked: { [key: ItemId]: Item } = {};
        const least: { [key: ItemId]: Item } = {};
        const items = state.groups[source.groupId].items;

        Object.keys(items).forEach((key) => {
          if (source.selected.includes(key)) {
            marked[key] = items[key];
          } else {
            least[key] = items[key];
          }
        });
        groups = {
          [source.groupId]: {
            label: "無題",
            state: "EDITING",
            items: least,
          },
          [action.newGroupId]: {
            label: "無題",
            state: "EDITING",
            items: marked,
          },
        };
      } else {
        groups = {
          [action.newGroupId]: {
            label: "無題",
            state: "EDITING",
            items: {},
          },
        };
      }

      return {
        ...state,
        groups: { ...state.groups, ...groups },
      };
    }
    case "UPLOAD_MANY":
      return {
        ...state,
        groups: mapDict(state.groups, (groupId, group) =>
          action.groupIds.includes(groupId)
            ? GroupReducer(group, {
              type: "UPLOAD",
              groupId,
            })
            : group
        ),
      };
    case "SET_GROUP_ITEMS":
    case "UPLOAD_COMPLETE":
      return {
        ...state,
        groups: mapDict(state.groups, (_, group) =>
          GroupReducer(group, action)
        ),
      };
    default:
      return state;
  }
};

export default StatusReducer;
