import { useCallback, useEffect, useMemo, useReducer, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { uuidv7 } from "uuidv7";
import reducer from "../app/reducer";
import "./UploadPage.css";
import GalleryView, { GalleryViewProps } from "./GalleryView";
import { ArchiveFiles, GenerateId } from "../app/libs";
import { Context } from "../app/context";
import { ApplicationError } from "../app/repositories/ApplicationError";
import { GroupId, Item, ItemId } from "../app/models";

function Never<T>(_: never[]): T {
  throw new Error("assert never")
}

type Props = {
  context: Context;
}

function UploadPage({ context }: Props) {
  const navigate = useNavigate();
  const location = useLocation();

  const DEFAULT_GROUP_ID = useMemo(() => GenerateId<GroupId>(), []);

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
  const [asyncError, setAsyncError] = useState<unknown>();

  useEffect(() => {
    const dict = Object.fromEntries(
      new URLSearchParams(location.search).entries()
    );
    setViewGroupId(dict["group"] ?? DEFAULT_GROUP_ID);
  }, [DEFAULT_GROUP_ID, location.search]);

  const onUpdateGroupItems = useCallback(
    (groupId: GroupId, items: { [key: ItemId]: Item }) => dispatch({
      type: "SET_GROUP_ITEMS",
      groupId: groupId,
      items: items,
    }),
    []
  );

  const onCreateGroup = useCallback(
    (groupId: GroupId, selectedItems: ItemId[]) => {
      dispatch({
        type: "CREATE_GROUP",
        newGroupId: GenerateId<GroupId>(),
        source: {
          groupId: groupId,
          selected: selectedItems,
        },
      });
    },
    []
  );

  const archiveAsync = useCallback(async (groupId: string, files: File[]) => {
    try {
      const key = await context.backend.upload(await ArchiveFiles(files));
      dispatch({
        type: "UPLOAD_COMPLETE",
        groupId: groupId,
        key,
      });
    } catch (err) {
      if (err instanceof ApplicationError) {
        dispatch({
          type: "UPLOAD_FAILED",
          groupId: groupId,
          error: err,
        });
      } else {
        setAsyncError(err);
      }
    }
  }, [])

  const onArchiveAll = useCallback(() => {
    const targetGroupIds = Object.keys(status.groups).filter((groupId) => {
      const g = status.groups[groupId];
      return g.state === "EDITING" && Object.keys(g.items).length > 0;
    });

    if (targetGroupIds.length === 0) {
      return
    }

    const removePid = (pid: string) => {
      setPromisePool(promisePool.filter((v) => v.pid !== pid));
    }

    const promiseList = targetGroupIds.map((groupId) => {
      const files = Object.values(status.groups[groupId].items).map(
        (item) => item.file
      );

      const pid = uuidv7();

      return {
        pid,
        promise: (async () => {
          await archiveAsync(groupId, files)
          removePid(pid);
        })(),
      };
    });

    setPromisePool(promisePool.concat(promiseList));
    dispatch({ type: "UPLOAD_MANY", groupIds: targetGroupIds });
  }, [status.groups]);

  const groupButtons = (
    <div style={{ display: "flex", flexDirection: "column" }}>
      {Object.keys(status.groups)
        .sort()
        .map((groupId, idx) => (
          <button key={idx} onClick={() => navigate(`?group=${groupId}`)}>
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
          dispatch({ type: "CREATE_GROUP", newGroupId: GenerateId() })
        }
      >
        +
      </button>
    </div>
  );

  const archiveAllButton = <button onClick={onArchiveAll}>Archive All</button>;

  if (!(viewGroupId && status.groups[viewGroupId])) {
    return null;
  }

  const group = status.groups[viewGroupId];
  const props: GalleryViewProps = group.state === "EDITING" ? {
    label: group.label,
    items: group.items,
    state: "EDITING",
    onCreateGroup: (ids: GroupId[]) => onCreateGroup(viewGroupId, ids),
    onUpdateItems: (items: { [key: ItemId]: Item }) => onUpdateGroupItems(viewGroupId, items),
  } : group.state === "ARCHIVING" ? {
    label: group.label,
    items: group.items,
    state: "ARCHIVING",
  } : group.state === "COMPLETE" ? {
    label: group.label,
    items: group.items,
    state: "COMPLETE",
    url: new URL(window.location.origin + '/download?key=' + group.key),
  } : group.state === "FAILED" ? {
    label: group.label,
    items: group.items,
    state: "FAILED",
    error: group.error,
  } : Never([group.state]);

  return (
    <div>
      <header>
        <p>Tsumiki Uploader</p>
      </header>
      <div className="body-container">
        <aside className="sidebar">
          <div className="sidebar-contents">{groupButtons}</div>
          <div>{archiveAllButton}</div>
        </aside>
        <section className="main-section">
          {asyncError ? <p>{`${asyncError}`}</p> : null}
          {(viewGroupId && status.groups[viewGroupId]) ? <GalleryView {...props} /> : null}
        </section>
      </div>
    </div>
  );
}

export default UploadPage;
