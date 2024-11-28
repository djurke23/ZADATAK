import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../services/auth.service'; // Uveri se da je putanja tačna
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';

@Component({
  standalone: true,
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css'],
  imports: [FormsModule, CommonModule],
})
export class LoginComponent {
  nickname: string = '';
  password: string = '';

  constructor(private authService: AuthService, private router: Router) {}

  onLogin() {
    if (!this.nickname || !this.password) {
      alert('Please enter both nickname and password.');
      return;
    }

    this.authService.login(this.nickname, this.password).subscribe({
      next: (response: { token: string }) => {
        localStorage.setItem('token', response.token); // Čuvanje tokena
        this.router.navigate(['/dashboard']); // Navigacija na Dashboard
      },
      error: (err: any) => {
        alert('Login failed: ' + (err.error?.message || 'Unknown error'));
      },
    });
  }
}
