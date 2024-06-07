
// User
export interface User {
  id: number;
  name: string;
  mail: string,
  gradeID: number,
  departmentID: number,
  bureauID: number;
  roleID: number;
  studentNumber: number,
  tel: string,
  password: string,
  createdAt?: string;
  updatedAt?: string;
}

// Grade(学年)
export interface Grade {
  id: number;
  grade: string;
  createdAt?: string;
  updatedAt?: string;
}

// Department(学科)
export interface Department {
  id: number;
  department: string;
  createdAt?: string;
  updatedAt?: string;
}

// Bureau(局)
export interface Bureau {
  id: number;
  bureau: string;
  createdAt?: string;
  updatedAt?: string;
}

// Weather(天気)
export interface Weather {
  id: number;
  weather: string;
  createdAt?: string;
  updatedAt?: string;
}

// Place(集合場所)
export interface Place {
  id: number;
  place: string;
  remark: string;
  createdAt?: string;
  updatedAt?: string;
}

// Task(仕事内容)
export interface Task {
  id: number;
  task: string;
  placeID: number;
  url: string;
  superviserID: number;
  color: string;
  remark: string;
  yearID: number;
  createdAt?: string;
  updatedAt?: string;
}

// Year(開催年度)
export interface Year {
  id: number;
  year: number;
  createdAt?: string;
  updatedAt?: string;
}

// Date(開催日)
export interface Date {
  id: number;
  yearID: number;
  name: string;
  date: string;
  createdAt?: string;
  updatedAt?: string;
}

// Time(時間)
export interface Time {
  id: number;
  time: string;
  createdAt?: string;
  updatedAt?: string;
}

// Shift(シフト)
export interface Shift {
  id: number;
  taskID: number;
  userID: number;
  yearID: number;
  dateID: number;
  timeID: number;
  weatherID: number;
  isAttendance: boolean;
  createdAt?: string;
  updatedAt?: string;
}