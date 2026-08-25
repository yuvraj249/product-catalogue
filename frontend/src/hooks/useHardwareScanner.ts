import { useEffect, useRef } from 'react';

interface HardwareScannerOptions {
  onScan: (barcode: string) => void;
  maxIntervalMs?: number;
  minBarcodeLength?: number;
}

export const useHardwareScanner = ({
  onScan,
  maxIntervalMs = 25,
  minBarcodeLength = 3,
}: HardwareScannerOptions) => {
  const bufferRef = useRef<string[]>([]);
  const lastKeyTimeRef = useRef<number>(0);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Avoid intercepting input fields unless designated with data-enable-scanner
      const target = e.target as HTMLElement;
      if (['INPUT', 'TEXTAREA', 'SELECT'].includes(target.tagName) && !target.dataset.enableScanner) {
        return;
      }

      const currentTime = performance.now();
      const timeDiff = currentTime - lastKeyTimeRef.current;
      lastKeyTimeRef.current = currentTime;

      if (e.key === 'Enter') {
        if (bufferRef.current.length >= minBarcodeLength) {
          e.preventDefault();
          const completeBarcode = bufferRef.current.join('');
          onScan(completeBarcode);
        }
        bufferRef.current = [];
        return;
      }

      // Scanner typing pace validation
      if (timeDiff > maxIntervalMs && bufferRef.current.length > 0) {
        bufferRef.current = []; // Wipe human typing noise
      }

      if (e.key.length === 1) { // Standard printable characters
        bufferRef.current.push(e.key);
      }
    };

    window.addEventListener('keydown', handleKeyDown, true);
    return () => window.removeEventListener('keydown', handleKeyDown, true);
  }, [onScan, maxIntervalMs, minBarcodeLength]);
};
