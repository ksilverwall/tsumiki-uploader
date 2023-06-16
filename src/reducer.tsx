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

const GroupReducer: Reducer<Group, GroupAction> = (state, action) => {
  switch (action.type) {
    case "UPLOAD":
      return {
        ...state,
        state: "ARCHIVING",
        promise: action.promise,
      };
    default:
      return state;
  }
};

const StatusReducer: Reducer<Status, Action> = (state, action) => {
  switch (action.type) {
    case "LOAD":
      return {
        ...state,
        groups: {
          ...state.groups,
          [action.groupId]: {
            ...state.groups[action.groupId],
            items: {
              ...state.groups[action.groupId].items,
              ...Object.fromEntries(
                action.items.map((item) => [
                  item.id,
                  {
                    file: item.file,
                    marked: false,
                  },
                ])
              ),
            },
          },
        },
      };
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
    case "UPLOAD":
      return {
        ...state,
        groups: mapDict(state.groups, (groupId, group) =>
          GroupReducer(group, {
            type: "UPLOAD",
            promise: action.promises[groupId],
          })
        ),
      };
    default:
      return state;
  }
};

export default StatusReducer;
