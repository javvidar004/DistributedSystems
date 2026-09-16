// components/shopping-list/WeekSelector.tsx
'use client';
import { ChevronLeftIcon, ChevronRightIcon } from '@heroicons/react/24/solid';
import React, { useEffect, useState } from 'react';

interface WeekSelectorProps {
  /** Called whenever the visible week changes. Receives startISO (YYYY-MM-DD) and endISO (YYYY-MM-DD) */
  onChange?: (startISO: string, endISO: string) => void;
  initialDate?: Date;
}

const formatShort = (d: Date) => d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });

const toISODate = (d: Date) => d.toISOString().slice(0, 10);

const getWeekRange = (date: Date) => {
  // compute Monday (start) and Sunday (end) for the week containing `date`
  const tmp = new Date(date);
  const day = tmp.getDay(); // 0 (Sun) .. 6 (Sat)
  const diffToMonday = (day + 6) % 7; // 0 if Monday, 6 if Sunday
  const start = new Date(tmp);
  start.setDate(tmp.getDate() - diffToMonday);
  start.setHours(0, 0, 0, 0);
  const end = new Date(start);
  end.setDate(start.getDate() + 6);
  end.setHours(23, 59, 59, 999);
  return { start, end };
};

/**
 * WeekSelector: computes a Monday-Sunday range and allows navigating weeks.
 * Calls `onChange(startISO, endISO)` when the week changes.
 */
const WeekSelector = ({ onChange, initialDate }: WeekSelectorProps) => {
  const today = initialDate ? new Date(initialDate) : new Date();
  const initialRange = getWeekRange(today);

  const [startDate, setStartDate] = useState<Date>(initialRange.start);

  useEffect(() => {
    // notify parent on mount and when startDate changes
    const end = new Date(startDate);
    end.setDate(startDate.getDate() + 6);
    if (onChange) onChange(toISODate(startDate), toISODate(end));
  }, [startDate]);

  const prevWeek = () => setStartDate(d => { const n = new Date(d); n.setDate(d.getDate() - 7); return n; });
  const nextWeek = () => setStartDate(d => { const n = new Date(d); n.setDate(d.getDate() + 7); return n; });
  const goToToday = () => setStartDate(getWeekRange(new Date()).start);

  const endDate = new Date(startDate);
  endDate.setDate(startDate.getDate() + 6);

  return (
    <div>
      <label className="block text-sm font-medium text-gray-700 mb-2">
        Selecciona el menú semanal:
      </label>
      <div className="flex items-center justify-between p-3 bg-gray-50 rounded-md border">
        <div className="flex items-center gap-2">
          <button onClick={prevWeek} className="p-2 rounded-full hover:bg-gray-200 transition-colors" aria-label="Previous week">
            <ChevronLeftIcon className="h-5 w-5 text-gray-600" />
          </button>
          <button onClick={goToToday} className="text-sm text-gray-600 underline">This week</button>
        </div>

        <span className="font-semibold text-dark-background">
          {`${formatShort(startDate)} - ${formatShort(endDate)}`}
        </span>

        <div className="flex items-center gap-2">
          <button onClick={nextWeek} className="p-2 rounded-full hover:bg-gray-200 transition-colors" aria-label="Next week">
            <ChevronRightIcon className="h-5 w-5 text-gray-600" />
          </button>
        </div>
      </div>
    </div>
  );
};

export default WeekSelector;