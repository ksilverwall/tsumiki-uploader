import { useCallback, useEffect, useMemo, useReducer, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { uuidv7 } from "uuidv7";
import reducer from "./reducer";
import "./App.css";
import GalleryView from "./GalleryView";

function generateId<T extends string>(): T {
  const newId = uuidv7();
  return newId as T;
}

function App() {
  const navigate = useNavigate();
  const location = useLocation();

  const DEFAULT_GROUP_ID = useMemo(() => generateId<GroupId>(), []);

  const [viewGroupId, setViewGroupId] = useState<string>(DEFAULT_GROUP_ID);
  const [status, dispatch] = useReducer(reducer, {
    groups: {
      [DEFAULT_GROUP_ID]: {
        label: "無題",
        state: "EDITING",
        items: {},
      },
    },
  });

  useEffect(() => {
    const dict = Object.fromEntries(
      new URLSearchParams(location.search).entries()
    );
    setViewGroupId(dict["group"] ?? DEFAULT_GROUP_ID);
  }, [DEFAULT_GROUP_ID, location.search]);

  const onLoad = useCallback(
    (groupId: GroupId, files: File[]) =>
      dispatch({
        type: "LOAD",
        groupId: groupId,
        items: files.map((f) => ({
          id: generateId<ItemId>(),
          file: f,
        })),
      }),
    []
  );

  const onCreateGroup = useCallback(
    (groupId: GroupId, selectedItems: ItemId[]) => {
      dispatch({
        type: "CREATE_GROUP",
        newGroupId: generateId<GroupId>(),
        source: {
          groupId: groupId,
          selected: selectedItems,
        },
      });
    },
    []
  );

  const groupButtons = (
    <div style={{ display: "flex", flexDirection: "column" }}>
      {Object.keys(status.groups)
        .sort()
        .map((v) => (
          <button onClick={() => navigate(`?group=${v}`)}>{v}</button>
        ))}
      <button
        onClick={() =>
          dispatch({ type: "CREATE_GROUP", newGroupId: generateId() })
        }
      >
        +
      </button>
    </div>
  );

  const archiveAllButton = <button>Archive All</button>;

  return (
    <>
      <header>
        <p>Tsumiki Uploader</p>
      </header>
      <div className="body-container">
        <aside className="sidebar">
          <div className="sidebar-contents">{groupButtons}</div>
          <div>{archiveAllButton}</div>
        </aside>
        <section className="main-section">
          <GalleryView
            groupId={viewGroupId}
            group={status.groups[viewGroupId]}
            onLoad={onLoad}
            onCreateGroup={onCreateGroup}
          />
        </section>
      </div>
    </>
  );
}

export default App;
