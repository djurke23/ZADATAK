import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class AuthService { // Preimenovano iz ApiService u AuthService
  private apiUrl = 'http://localhost:8080'; // URL backend-a

  constructor(private http: HttpClient) {}

  // Login metoda
  login(nickname: string, password: string): Observable<any> {
    return this.http.post(`${this.apiUrl}/login`, { nickname, password });
  }

  // Dobijanje svih korisnika
  getUsers(): Observable<any> {
    const headers = { Authorization: `Bearer ${localStorage.getItem('token')}` };
    return this.http.get(`${this.apiUrl}/users/`, { headers });
  }  

  // Dodavanje korisnika
  addUser(user: any): Observable<any> {
    const headers = { Authorization: `Bearer ${localStorage.getItem('token')}` };
    return this.http.post(`${this.apiUrl}/users/`, user, { headers });
  }

  // Ažuriranje korisnika
  updateUser(id: number, user: any): Observable<any> {
    const headers = { Authorization: `Bearer ${localStorage.getItem('token')}` };
    return this.http.put(`${this.apiUrl}/users/${id}`, user, { headers });
  }

  // Brisanje korisnika
  deleteUser(id: number): Observable<any> {
    const headers = { Authorization: `Bearer ${localStorage.getItem('token')}` };
    return this.http.delete(`${this.apiUrl}/users/${id}`, { headers });
  }
}
