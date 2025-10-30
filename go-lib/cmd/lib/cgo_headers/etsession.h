// Copyright (c) 2024 Proton AG
//
// This file is part of Proton Mail Bridge.
//
// Proton Export Tool is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Proton Export Tool is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Proton Mail Bridge. If not, see <https://www.gnu.org/licenses/>.

#ifndef ET_SESSION_H
#define ET_SESSION_H

#include <stdint.h>
#include <stdlib.h>

typedef const char cchar_t;

typedef struct etSession etSession;

typedef enum etSessionStatus {
	ET_SESSION_STATUS_OK,
	ET_SESSION_STATUS_ERROR,
	ET_SESSION_STATUS_INVALID,
	ET_SESSION_STATUS_CANCELLED,
} etSessionStatus;

typedef enum etSessionLoginState {
	ET_SESSION_LOGIN_STATE_LOGGED_OUT,
	ET_SESSION_LOGIN_STATE_AWAITING_TOTP,
	ET_SESSION_LOGIN_STATE_AWAITING_HV,
	ET_SESSION_LOGIN_STATE_AWAITING_MAILBOX_PASSWORD,
	ET_SESSION_LOGIN_STATE_LOGGED_IN,
} etSessionLoginState;

typedef struct etSessionCallbacks {
    void *ptr;
    void (*onNetworkLost)(void*);
    void (*onNetworkRestored)(void*);
} etSessionCallbacks;

#endif // ET_SESSION_H