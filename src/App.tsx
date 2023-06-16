import { useCallback, useEffect, useMemo, useReducer, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { uuidv7 } from "uuidv7";
import reducer from "./reducer";
import "./App.css";
import GalleryView from "./GalleryView";
import Uploader from "./Uploader";

function generateId<T extends string>(): T {
  const newId = uuidv7();
  return newId as T;
}

function App() {
  const navigate = useNavigate();
  const location = useLocation();

  const DEFAULT_GROUP_ID = useMemo(() => generateId<GroupId>(), []);

  const [promisePool, setPromisePool] = useState<
    { pid: string; promise: Promise<void> }[]
  >([]);
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

  const onArchiveAll = () => {
    const targetGroupIds = Object.keys(status.groups).filter((groupId) => {
      const g = status.groups[groupId];
      return g.state === "EDITING" && Object.keys(g.items).length > 0;
    });

    const promiseList = targetGroupIds.map((groupId) => {
      const files = Object.values(status.groups[groupId].items).map(
        (item) => item.file
      );

      const pid = uuidv7();

      return {
        pid,
        promise: (async () => {
          try {
            const key = await new Uploader().upload(files);
            dispatch({
              type: "UPLOAD_COMPLETE",
              groupId: groupId,
              key,
            });
          } catch (err) {
            // TODO: Store error
            console.log(err);
          }

          setPromisePool(promisePool.filter((v) => v.pid !== pid));
        })(),
      };
    });

    setPromisePool(promisePool.concat(promiseList));
    dispatch({ type: "UPLOAD_MANY", groupIds: targetGroupIds });
  };

  const groupButtons = (
    <div style={{ display: "flex", flexDirection: "column" }}>
      {Object.keys(status.groups)
        .sort()
        .map((groupId) => (
          <button onClick={() => navigate(`?group=${groupId}`)}>
            <div style={{ display: "flex", flexDirection: "row" }}>
              <div>
                <p>{`[${status.groups[groupId].state}]`}</p>
              </div>
              <div>
                <p>{groupId}</p>
              </div>
            </div>
          </button>
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

  const archiveAllButton = <button onClick={onArchiveAll}>Archive All</button>;

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
