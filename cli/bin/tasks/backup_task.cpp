// Copyright (c) 2023 Proton AG
//
// This file is part of Proton Export Tool.
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
// along with Proton Export Tool.  If not, see <https://www.gnu.org/licenses/>.

#include "tasks/backup_task.hpp"
#include <etsession.hpp>
#include <iostream>

BackupTask::BackupTask(etcpp::Session& session, const std::filesystem::path& backupPath, const FilterOptions& filterOptions) :
    mBackup(session.newBackup(
        backupPath.u8string().c_str(),
        filterOptions.labelIDs.c_str(),
        filterOptions.sender.c_str(),
        filterOptions.recipient.c_str(),
        filterOptions.domain.c_str(),
        filterOptions.after.c_str(),
        filterOptions.before.c_str(),
        filterOptions.subject.c_str()
    )) {}

// Backward compatibility constructor
BackupTask::BackupTask(etcpp::Session& session, const std::filesystem::path& backupPath, const char* labelIDs) :
    mBackup(session.newBackup(backupPath.u8string().c_str(), labelIDs)) {}

void BackupTask::onProgress(float progress) {
    updateProgress(progress);
}

void BackupTask::run() {
    mBackup.start(*this);
}

void BackupTask::cancel() {
    mBackup.cancel();
}

std::string_view BackupTask::description() const {
    return "Export Mail";
}
