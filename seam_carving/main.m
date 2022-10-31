% NOTE
%--------------------------------------------------------------------------
% This code has been divided into the following parts:
%--------------------------------------------------------------------------
% Main body: 
% The first section. Makes a function call to
% compute_energy_image to get the energy image based on the L2 norm energy
% function and to apply_seam_carve, which is the helper function allowing
% us to handle either horizontal or vertical seams.
%--------------------------------------------------------------------------
% compute_energy_image: 
% This function handles image energy computation via
% L2 norm.
%--------------------------------------------------------------------------
% apply_seam_carve: 
% This function is a helper function that computes the
% cumulative_energy_map by calling the respective function
% and thus finds and applies either the horizontal or
% vertical seam, with the use of the find_horizontal_seam and
% find_vertical_seam functions.
%--------------------------------------------------------------------------
% cumulative_energy_map:
% This function applies dynamic programming techniques on the computed
% energy image and finds out the required seam's opt array (in this case,
% named cumulative_energy_map).
%--------------------------------------------------------------------------
% find_horizontal/vertical_seam:
% These two functions are responsible for providing the indices of the
% 4-connected vertical or horizontal seam pixels that have the least energy.
%--------------------------------------------------------------------------


% Main body code begins here
%--------------------------------------------------------------------------
% These three variables can be freely modified to change the number of
% seams, their types, and the filename of the image to be used.
seam_type = 'vertical';
number_of_seams = 100;
filename = '00000005';

input_image=imread(append('Input/archive/',filename,'.jpg'));
whos input_image

% The energy_image function is defined below.
% It utilizes an energy function to evaluate the energy posessed by each
% pixel, which is the first step to seam carving.
computed_energy = compute_energy_image(input_image);

carved_image=input_image;
modified_energy_image=computed_energy;

% Depending on the type of the seam, apply either decrease_height or
% decrease_width, providing the current state of the input image and the
% energy image as arguments.

for i=1:number_of_seams
    [carved_image,modified_energy_image] = apply_seam_carve(carved_image,modified_energy_image,seam_type);
end


% Print the output and show the results in a side by side comparison
output_path = append('Output/output_',seam_type,'_carving_',filename,'.png');
imwrite(carved_image,output_path);

subplot(1,2,1),imshow(input_image),title(append('Original ',filename));
subplot(1,2,2),imshow(carved_image),title(append(num2str(number_of_seams),'px ',seam_type,' seams carved ',filename));


% -------------------
% MAIN BODY ENDS HERE
%--------------------


% Acquiring energy image
function calculated_energy = compute_energy_image(input_image)
    % Energy function works on magnitude, so grayscale conversion is done
    % to allow it to work smoothly.
    grayscale_conversion=double(rgb2gray(input_image));
    
    % This is a direct implementation of the energy function in the
    % original literature. Partial differentiation is simply the
    % subtraction of an element of a certain index from its predecessor on
    % the x or y axis.
    dI_by_dx=grayscale_conversion(:,2:end)-grayscale_conversion(:,1:end-1);
    dI_by_dy=grayscale_conversion(2:end,:)-grayscale_conversion(1:end-1,:);
    
    % The L2 norm here is used as an energy function. It is defined as the
    % square root of the sum of squares of gradients.
    calculated_energy= sqrt(dI_by_dx(1:end-1,:).^2+dI_by_dy(:,1:end-1).^2);
    calculated_energy(end+1,end+1)=0;
end


% -------------------------
% ENERGY FUNCTION ENDS HERE
%--------------------------


% A helper function to execute seam carving
function [carved_image,modified_energy_image] = apply_seam_carve(input_image,energy_image,seam_type)
    cumulative_energy_map = cumulative_min_energy_map(energy_image,seam_type);
    % This function has a demonstration built into it. 
    % We color the seams red and show a pseudo-animation of the carve as it
    % is in progress.
    original_image = input_image;
    demonstration_image = input_image;
    if strcmp('horizontal',seam_type)
        seam = find_horizontal_seam(cumulative_energy_map);
        for seam_index=1:length(seam)
            demonstration_image(seam(seam_index),seam_index,1) = 255;
            demonstration_image(seam(seam_index),seam_index,2) = 0;
            demonstration_image(seam(seam_index),seam_index,3) = 0;
            input_image(seam(seam_index):end-1,seam_index,:)=input_image(seam(seam_index)+1:end,seam_index,:);
        end
        carved_image=input_image(1:end-1,:,:);
    else
        seam = find_vertical_seam(cumulative_energy_map);
        for seam_index=1:length(seam)
            demonstration_image(seam_index,seam(seam_index),1) = 255;
            demonstration_image(seam_index,seam(seam_index),2) = 0;
            demonstration_image(seam_index,seam(seam_index),3) = 0;
            input_image(seam_index,seam(seam_index):end-1,:)=input_image(seam_index,seam(seam_index)+1:end,:);
        end
        carved_image=input_image(:,1:end-1,:);
    end
    imshow(original_image)
    imshow(demonstration_image)
    imshow(carved_image)
    modified_energy_image=compute_energy_image(carved_image);
end


% ------------------------------------
% SEAM CARVE HELPER FUNCTION ENDS HERE
%-------------------------------------


% Finding the required cumulative energy map via dynamic programming.
function cumulative_energy_map = cumulative_min_energy_map(energy_image,seam_type)
    output_map=zeros(size(energy_image));
    if strcmp('horizontal',seam_type)   
        output_map(:,1)=energy_image(:,1);
        output_map(1,2:end)=inf;
        output_map(end,2:end)=inf;
        for col=2:size(output_map,2)
            for row=2:size(output_map,1)-1
                output_map(row,col)=energy_image(row,col)+min(output_map(row-1:row+1,col-1));
            end
        end
    else
        % As described in the paper, we set the first row of the energy map
        % as being the same as the original computed energy image
        output_map(1,:)=energy_image(1,:);
        
        % The output_map is like an opt[] array in a dynamic programming
        % algorithm; All other entries are set to infinity since this is a
        % minimization problem.
        output_map(2:end,1)=inf;
        output_map(2:end,end)=inf;
        for row=2:size(output_map,1)
            for col=2:size(output_map,2)-1
                output_map(row,col)=energy_image(row,col)+min(output_map(row-1,col-1:col+1));
            end
        end
    end
    cumulative_energy_map=output_map;
end


% ---------------------------------------------------
% FUNCTION TO COMPUTE CUMULATIVE ENERGY MAP ENDS HERE
%----------------------------------------------------


% Helper functions to find and return the indices of the horizontal and
% vertical seams.
function horizontal_seam = find_horizontal_seam(cumulative_energy_map)
    horizontal_seam=zeros(size(cumulative_energy_map,2),1);
    [min_energy_value,min_energy_index]=min(cumulative_energy_map(:,end));
    horizontal_seam(end)=min_energy_index;
    for seam_index=size(cumulative_energy_map,2)-1:-1:1
        % This is a optimal path retrieval of the dynamic programming
        % algorithm that gives us the cumulative energy map
        [min_energy_value,local_min_energy_index]=min(cumulative_energy_map(min_energy_index-1:min_energy_index+1,seam_index));
        
        % local_min_energy_index is the 3-element column array that we find, 
        % we need to index it to the correct indexing posessed by the image,
        % hence the min_energy_index+local_min_energy_index-2
        min_energy_index=min_energy_index+local_min_energy_index-2;
        horizontal_seam(seam_index)=min_energy_index;
    end
end

function vertical_seam = find_vertical_seam(cumulative_energy_map)
    vertical_seam=zeros(size(cumulative_energy_map,1),1);
    [min_energy_value,min_energy_index]=min(cumulative_energy_map(end,:));
    vertical_seam(end)=min_energy_index;
    for seam_index=size(cumulative_energy_map,1)-1:-1:1
        %cumulative_energy_map(seam_index,min_energy_index-1:min_energy_index+1)
        % This is a optimal path retrieval of the dynamic programming
        % algorithm that gives us the cumulative energy map
        [min_energy_value,local_min_energy_index]=min(cumulative_energy_map(seam_index,min_energy_index-1:min_energy_index+1));
        
        % local_min_energy_index is the 3-element row array that we find, 
        % we need to index it to the correct indexing posessed by the image,
        % hence the min_energy_index+local_min_energy_index-2
        min_energy_index=min_energy_index+local_min_energy_index-2;
        vertical_seam(seam_index)=min_energy_index;
    end
end


% -------------------------------------------------
% HORIZ/VERT SEAM FINDING HELPER FUNCTIONS END HERE
%--------------------------------------------------