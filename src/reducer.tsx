const StatusReducer = (state: Status, action: Action): Status => {
  switch (action.type) {
    case "LOAD":
      return {
        ...state,
        groups: {
          ...state.groups,
          [action.groupId]: {
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
      let groups;
      const source = action.source;
      if (source) {
        const marked: {[key: ItemId]: Item} = {};
        const least: {[key: ItemId]: Item} = {};
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
            items: least,
          },
          [action.newGroupId]: {
            items: marked,
          },
        };
      } else {
        groups = {
          [action.newGroupId]: {
            items: {},
          },
        };
      }

      return {
        ...state,
        groups: { ...state.groups, ...groups },
      };
    }
    default:
      return state;
  }
};

export default StatusReducer;
