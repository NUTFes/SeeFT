
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