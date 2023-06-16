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
    case "CREATE_GROUP":
      return {
        ...state,
        groups: {
          ...state.groups,
          [action.newGroupId]: {
            items: {},
          },
        },
      };
    case "ARCHIVE": {
      const marked: any = {};
      const least: any = {};
      Object.keys(state.groups[action.groupId].items).forEach((key) => {
        if (action.source.selected.includes(key)) {
          marked[key] = state.groups[action.groupId].items[key];
        } else {
          least[key] = state.groups[action.groupId].items[key];
        }
      });
      return {
        ...state,
        groups: {
          ...state.groups,
          [action.groupId]: {
            items: least,
          },
          [action.newGroupId]: {
            items: marked,
          },
        },
      };
    }
    default:
      return state;
  }
};

export default StatusReducer;
