import "./PreloadOverlay.css";

interface PreloadOverlayProps {
    type?: "spinner" | "dots"
}

export const PreloadOverlay = ({type="spinner"}:PreloadOverlayProps) => {
    return (
        <div className="preload-overlay">
            <span className={`${type}-loader`}></span>
        </div>
    )
}
