interface GameControlsProps {
  onResign: () => void;
  onNewGame: () => void;
  disabled: boolean;
  gameOver: boolean;
}

export function GameControls({ onResign, onNewGame, disabled, gameOver }: GameControlsProps) {
  const handleResign = () => {
    if (window.confirm("Resign this game? You will lose.")) {
      onResign();
    }
  };

  return (
    <div className="game-controls">
      <button
        type="button"
        onClick={handleResign}
        disabled={disabled || gameOver}
        className="control-btn resign-btn"
      >
        Resign
      </button>
      <button
        type="button"
        onClick={onNewGame}
        className="control-btn new-game-btn"
      >
        New Game
      </button>
    </div>
  );
}
