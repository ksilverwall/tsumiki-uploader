const GalleryViewLayout: React.FC<{ slots: { header: JSX.Element, images: JSX.Element[] } }> = ({ slots }) => {
    return (
        <div className="gallery-view">
            <div>{slots.header}</div>
            <div className="image-list">
                {slots.images.map((e, idx) => <div key={idx}>{e}</div>)}
            </div>
        </div>
    );
}
export default GalleryViewLayout
